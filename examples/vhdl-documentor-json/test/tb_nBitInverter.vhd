
library IEEE;
use IEEE.STD_LOGIC_1164.all;

entity NBitInverter_tb is
end NBitInverter_tb;

architecture Behavioral of NBitInverter_tb is

    constant N : integer := 4;

    signal Input  : std_logic_vector(N-1 downto 0) := (others => '0');
    signal Output : std_logic_vector(N-1 downto 0);

begin

    uut : entity work.NBitInverter
        generic map (N => N)
        port map (
            Input  => Input,
            Output => Output
            );

    -- Stimulus process to apply test vectors
    stim_proc : process
    begin
        -- Test Case 1: All zeros
        Input <= "0000";
        wait for 10 ns;

        -- Test Case 2: All ones
        Input <= "1111";
        wait for 10 ns;

        -- Test Case 3: Alternating bits (1010)
        Input <= "1010";
        wait for 10 ns;

        -- Test Case 4: Alternating bits (0101)
        Input <= "0101";
        wait for 10 ns;

        -- Test Case 5: Random value
        Input <= "1100";
        wait for 10 ns;

        -- Finish simulation
        wait;
    end process;

end Behavioral;
