library IEEE;
use IEEE.std_logic_1164.all;

entity tb_NMux2t1 is end tb_NMux2t1;

architecture behavior of tb_NMux2t1 is
    component nmux2t1
        port (
            AMux   : in  std_logic;
            BMux   : in  std_logic;
            Sel    : in  std_logic;
            Output : out std_logic
            );


    end component;

    --Inputs
    signal s_AMux : std_logic := '0';
    signal s_BMux : std_logic := '0';
    signal s_Out  : std_logic := '0';

    --Outputs
    signal f : std_logic;

begin
    DUT0 : nmux2t1
        port map (
            AMux   => s_AMux,
            BMux   => s_BMux,
            Sel    => s_Out,
            Output => f
            );
    -- Stimulus process
    stim_proc : process
    begin
        s_AMux <= '0';
        s_BMux <= '0';
        s_Out  <= '0';
        wait for 100 ns;

        assert f = '0' report "Test failed for s=0" severity error;

        wait;
    end process stim_proc;
end behavior;
