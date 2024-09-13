library IEEE;

use IEEE.STD_LOGIC_1164.all;
use IEEE.NUMERIC_STD.all;
use work.FullAdder;

entity NBitAdder is
    generic (
        N : integer := 4
        );
    port (
        A        : in  std_logic_vector(N-1 downto 0);
        B        : in  std_logic_vector(N-1 downto 0);
        Sum      : out std_logic_vector(N-1 downto 0);
        CarryOut : out std_logic_vector(N-1 downto 0)
        );
end NBitAdder;


architecture Behavioral of NBitAdder is
    signal carries : std_logic_vector(N downto 0);
begin
    -- Initialize the carry-in for the first adder to 0
    carries(0) <= '0';

    gen_full_adders : for i in 0 to N-1 generate
        full_adder_inst : entity FullAdder
            port map (
                i_A  => A(i),
                i_B  => B(i),
                Cin  => carries(i),
                Sum  => Sum(i),
                Cout => carries(i+1)
                );
    end generate;

    -- The carry-out of the last full adder is the CarryOut of the N-bit adder
    CarryOut <= carries(N-1 downto 0);

end Behavioral;
